#include "cal_1.h"
#include "ui_cal_1.h"

cal_1::cal_1(QString uid, QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::cal_1), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);

    ui->events_scroll->setWidgetResizable(true);
    loadEventsForToday();
}

cal_1::~cal_1()
{
    delete ui;
}


void cal_1::loadEventsForToday()
{
    QWidget* prev = ui->events_scroll->widget();
    if (prev != nullptr) {
        delete prev;
    }

    QWidget *central = new QWidget;
    QVBoxLayout* layout = new QVBoxLayout(central);

    QDate today = QDate::currentDate();
    QVector<EventData> events = cli.getEventsByDay(today, uid);


    if (events.isEmpty()) {
        QLabel *noEventsLabel = new QLabel("На сегодня событий нет");
        noEventsLabel->setAlignment(Qt::AlignCenter);
        layout->addWidget(noEventsLabel);
    } else {

        for (const auto &e : events) {
            event_entry *item = new event_entry(e, uid);
            connect(item, &event_entry::deleted, this, &cal_1::eventDeleted);
            layout->addWidget(item);
        }
    }

    central->setLayout(layout);
    ui->events_scroll->setWidget(central);
}

void cal_1::eventDeleted(event_entry* ev)
{
    bool ok = cli.deleteEvent(uid, ev->getData().ID);
    if (!ok) {
        QMessageBox::warning(this, "Ошибка", "При удалении события произошла ошибка");
    } else {
        delete ev;
    }

    loadEventsForToday();
}

void cal_1::showEvent(QShowEvent *event)
{
    QWidget::showEvent(event);
    loadEventsForToday();
}
